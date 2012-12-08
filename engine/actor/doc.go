/*
Package actor represents Actors, the base of all objects and characters in
gogame. Actors may recieve Messages and broadcast Messages to Subscribers.

In this documentation, UP is defined as the direction toward TopLevel and
DOWN is defined as the direction toward Actor.

Initialization moves DOWN, with each type either standing as a middleman
for a given channel or simply handing the channel down the line.

Messages move UP, with each type unable to handle the Message simply passing
it to the next channel. If a Message reaches the TopLevel, a panic occurs.

Broadcasts move DOWN. Generally, they move directly to the Actor and are
sent to all the currently registered Subscribers with the same Kind.
*/
package actor
